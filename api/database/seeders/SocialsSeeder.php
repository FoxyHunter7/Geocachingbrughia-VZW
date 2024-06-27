<?php

namespace Database\Seeders;

use App\Models\Social;
use Illuminate\Database\Seeder;

class SocialsSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/socials.csv'), ';');

        $model = new Social();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
