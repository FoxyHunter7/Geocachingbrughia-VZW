<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class StaticSiteProperty extends Model
{
    use HasFactory;

    protected $primaryKey = "property";
    protected $keyType = 'string';
    public $incrementing = false;

    protected $hidden = [
        'created_at',
        'updated_at'
    ];

    protected $fillable = [
        'property'
    ];

    public function contents()
    {
        return $this->hasMany(staticSiteContent::class);
    }
}
